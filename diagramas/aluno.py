# coding=System
from timestamp import *

class aluno(object):

  """
   

  :version:
  :author:
  """

  """ ATTRIBUTES

   

  id  (public)

   

  nome  (public)

   

  curso  (public)

   

  ano  (public)

   

  unidade  (public)

   

  rgm  (public)

   

  senha  (public)

   

  createdAt  (public)

   

  updatedAt  (public)

  """

  def getAll(self, db):
    """
     

    @param *sql.DB db : 
    @return ([]*Aluno,error) :
    @author
    """
    func (this Aluno) GetAll(db *sql.DB) ([]*Aluno, error) {
        alunos := make([]*Aluno, 0)
        selDB, err := db.Query("Select * from cadastros.alunos")
        if err != nil {
            return alunos, err
        }
        for selDB.Next() {
            auxAluno := new(Aluno)
            err = selDB.Scan(
                &auxAluno.ID,
                &auxAluno.Nome,
                &auxAluno.Curso,
                &auxAluno.Ano,
                &auxAluno.Unidade,
                &auxAluno.Rgm,
                &auxAluno.Senha,
                &auxAluno.CreatedAt,
                &auxAluno.UpdatedAt,
            )
            if err != nil {
                return alunos, err
            }
            alunos = append(alunos, auxAluno)
        }
        return alunos, err
    }




